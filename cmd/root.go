// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"storj.io/common/base58"
	"storj.io/common/encryption"
	"storj.io/common/grant"
	"storj.io/common/identity"
	"storj.io/common/macaroon"
	"storj.io/common/pb"
	"storj.io/common/peertls/tlsopts"
	"storj.io/common/rpc"
	"storj.io/common/storj"
	"storj.io/uplink"
	"storj.io/uplink/private/metaclient"
)

// variable for output-empty-passphrase-access flag
var outputtingGrant bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "accesschecker",
	Short: "Check for empty passphrase in access grant",
	Long:  `TBD`,
	Run: func(cmd *cobra.Command, args []string) {
		err := beginCheck(args[0])
		if err != nil {
			panic(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&outputtingGrant, "output-empty-passphrase-access", "o", false, "Output generated empty passphrase access")
}

func beginCheck(accessArg string) error {
	ctx := context.Background()

	scope, err := parseAccessRaw(accessArg)
	if err != nil {
		return errs.New("error message %+v", err)
	}

	apiKey, err := macaroon.ParseRawAPIKey(scope.ApiKey)
	if err != nil {
		return errs.New("error message %+v", err)
	}

	ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{
		Difficulty:  0,
		Concurrency: 1,
	})
	if err != nil {
		return errs.New("error message %+v", err)
	}

	tlsConfig := tlsopts.Config{
		UsePeerCAWhitelist: false,
		PeerIDVersions:     "0",
	}

	tlsOptions, err := tlsopts.NewOptions(ident, tlsConfig, nil)
	if err != nil {
		return errs.New("error message %+v", err)
	}

	dialer := rpc.NewDefaultDialer(tlsOptions)
	metainfo, err := metaclient.DialNodeURL(ctx, dialer, scope.SatelliteAddr, apiKey, "")
	if err != nil {
		return errs.New("error message %+v", err)
	}
	defer func() {
		err = metainfo.Close()
		if err != nil {
			fmt.Printf("error closing metainfo dialer %w\n", err)
		}
	}()

	info, err := metainfo.GetProjectInfo(ctx)
	if err != nil {
		return errs.New("error message %+v", err)
	}

	accessGrant, err := genAccessGrantEmptyPass(scope.SatelliteAddr, apiKey.Serialize(), info.ProjectSalt)
	if err != nil {
		return errs.New("error message %+v", err)
	}
	keys, n, err := checkFiles(ctx, accessGrant)
	if err != nil {
		return errs.New("error message %+v", err)
	}
	if n > 0 {
		if outputtingGrant {
			fmt.Printf("Generated empty passphrase grant: %s\n", accessGrant)
		}
		fmt.Printf("You have %d files uploaded without encryption using this api key:\n", n)
		for _, k := range keys {
			fmt.Println(k)
		}
		return nil
	}
	fmt.Println("You do not have any files uploaded without encryption using this api key")
	return nil
}

func checkFiles(ctx context.Context, access string) (keys []string, n int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	accessGrant, err := uplink.ParseAccess(access)
	if err != nil {
		return keys, 0, errs.Wrap(err)
	}
	project, err := uplink.OpenProject(ctx, accessGrant)
	if err != nil {
		return keys, 0, errs.Wrap(err)
	}
	defer func(project *uplink.Project) {
		closeErr := project.Close()
		if closeErr != nil {
			err = errs.Combine(err, closeErr)
		}
	}(project)

	buckets := project.ListBuckets(ctx, nil)

	for buckets.Next() {
		if buckets.Err() != nil {
			return keys, 0, errs.Wrap(buckets.Err())
		}
		b := buckets.Item()
		objects := project.ListObjects(ctx, b.Name, nil)
		for objects.Next() {
			if objects.Err() != nil {
				return keys, 0, errs.Wrap(objects.Err())
			}
			fullKey := "sj://" + b.Name + "/" + objects.Item().Key
			keys = append(keys, fullKey)
			n++
		}
	}

	return keys, n, nil

}

func genAccessGrantEmptyPass(satelliteNodeURL, apiKey string, projectSalt []byte) (string, error) {
	parsedAPIKey, err := macaroon.ParseAPIKey(apiKey)
	if err != nil {
		return "", err
	}

	const concurrency = 8
	key, err := encryption.DeriveRootKey([]byte(""), projectSalt, "", concurrency)

	encAccess := grant.NewEncryptionAccessWithDefaultKey(key)
	encAccess.SetDefaultPathCipher(storj.EncAESGCM)
	encAccess.LimitTo(parsedAPIKey)

	accessString, err := (&grant.Access{
		SatelliteAddress: satelliteNodeURL,
		APIKey:           parsedAPIKey,
		EncAccess:        encAccess,
	}).Serialize()
	if err != nil {
		return "", err
	}
	return accessString, nil

}

// parseAccessRaw decodes Scope from base58 string, that contains SatelliteAddress, ApiKey, and EncryptionAccess.
func parseAccessRaw(access string) (_ *pb.Scope, err error) {
	data, version, err := base58.CheckDecode(access)
	if err != nil || version != 0 {
		return nil, errs.New("invalid access grant format: %w", err)
	}

	p := new(pb.Scope)
	if err := pb.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil
}
