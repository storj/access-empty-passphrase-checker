// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"context"
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
	Use:   "accesschecker [accessgrant]",
	Short: "Check for files uploaded with an empty passphrase in a project",
	Long: `This command is used to determine if you have previously uploaded unencrypted files to a Storj DCS project.
			To use, go to the Satellite UI, and generate a new access grant for the project you are interested in.
			Then, run the command with the access you copied from the Satellite UI. It will tell you if you have any unencrypted files uploaded, and what they are called.
			If you need an (unsafe) "no passphrase" access so that you can download and remove your unencrypted files using uplink, run with the flag "--output-empty-passphrase-access".`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("You are missing an argument.\n")
			cmd.Println(cmd.UsageString())
			return
		}
		err := beginCheck(cmd, args[0])
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
	rootCmd.Flags().BoolVarP(&outputtingGrant, "output-empty-passphrase-access", "o", false, "Output generated empty passphrase access (only use to download/remove files)")
}

func beginCheck(cmd *cobra.Command, accessArg string) error {
	ctx := context.Background()

	scope, err := parseAccessRaw(accessArg)
	if err != nil {
		return errs.New("Could not parse access: %+v", err)
	}

	apiKey, err := macaroon.ParseRawAPIKey(scope.ApiKey)
	if err != nil {
		return errs.New("Could not parse API key from access: %+v", err)
	}

	ident, err := identity.NewFullIdentity(ctx, identity.NewCAOptions{
		Difficulty:  0,
		Concurrency: 1,
	})
	if err != nil {
		return errs.New("Could not create identity to dial satellite: %+v", err)
	}

	tlsConfig := tlsopts.Config{
		UsePeerCAWhitelist: false,
		PeerIDVersions:     "0",
	}

	tlsOptions, err := tlsopts.NewOptions(ident, tlsConfig, nil)
	if err != nil {
		return errs.New("Could not create TLS options: %+v", err)
	}

	dialer := rpc.NewDefaultDialer(tlsOptions)
	metainfo, err := metaclient.DialNodeURL(ctx, dialer, scope.SatelliteAddr, apiKey, "")
	if err != nil {
		return errs.New("Error dialing satellite metainfo: %+v", err)
	}
	defer func() {
		err = metainfo.Close()
		if err != nil {
			cmd.Printf("error closing metainfo dialer %w\n", err)
		}
	}()

	info, err := metainfo.GetProjectInfo(ctx)
	if err != nil {
		return errs.New("Error getting project info %+v", err)
	}

	accessGrant, err := genAccessGrantEmptyPass(scope.SatelliteAddr, apiKey.Serialize(), info.ProjectSalt)
	if err != nil {
		return errs.New("Error creating an 'empty passphrase' access %+v", err)
	}
	keys, n, err := checkFiles(ctx, accessGrant)
	if err != nil {
		return errs.New("Error checking for unencrypted files: %+v", err)
	}
	if n > 0 {
		cmd.Printf("\nYou have %d files uploaded without encryption in this project:\n\n", n)
		for _, k := range keys {
			cmd.Println(k)
		}

		if outputtingGrant {
			cmd.Println("\n==================================================\n")
			cmd.Printf("Generated empty passphrase grant:\n\n%s\n\n", accessGrant)
			cmd.Println("WARNING: This access is capable of uploading and downloading unencrypted files to and from your project. We recommend using it only to download and subsequently remove files which are unencrypted.")
		}
		return nil
	}
	cmd.Println("You do not have any files uploaded without encryption in this project.")
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
		if objects.Err() != nil {
			return keys, 0, errs.Wrap(objects.Err())
		}
	}
	if buckets.Err() != nil {
		return keys, 0, errs.Wrap(buckets.Err())
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
