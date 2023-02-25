package awssync

import (
	"context"
	appConfig "internal/config"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

const (
	PROFILE     string = "profile"
	ACCESS_KEY  string = "accessKey"
	IAM_PROFILE string = "iamProfile"
)

// LoadAwsConfig loads the AWS configuration
func loadAwsConfig(ctx context.Context, c appConfig.Config) aws.Config {
	switch c.Aws.Method {
	case PROFILE:
		// Load from profile
		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithSharedConfigProfile(c.Aws.Profile),
			config.WithRegion(c.Aws.Region),
		)
		if err != nil {
			panic(err)
		}
		return cfg
	case ACCESS_KEY:
		// Load from access key
		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithRegion(c.Aws.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(c.Aws.AccessKey, c.Aws.SecretKey, ""),
			),
		)
		if err != nil {
			panic(err)
		}
		return cfg
	case IAM_PROFILE:
		// Load from EC2 Instance IAM profile
		cfg, err := config.LoadDefaultConfig(
			ctx,
			config.WithRegion(c.Aws.Region),
			config.WithEC2IMDSEndpoint(""),
		)
		if err != nil {
			panic(err)
		}
		return cfg
	}
	panic("Invalid AWS method")
}

// GetIAMClient returns a new IAM client
func getIAMClient(cfg aws.Config) *iam.Client {
	iamClient := iam.NewFromConfig(cfg)
	return iamClient
}

// GetKeysForUser returns a list of SSH keys for a user
func getKeysForUser(ctx context.Context, c *iam.Client, u []string) ([]string, error) {
	var keys []string
	for _, user := range u {
		// Get SSH keys for user
		resp, err := c.ListSSHPublicKeys(
			ctx,
			&iam.ListSSHPublicKeysInput{UserName: &user},
		)
		if err != nil {
			return nil, err
		}
		for _, key := range resp.SSHPublicKeys {
			resp, err := c.GetSSHPublicKey(
				ctx,
				&iam.GetSSHPublicKeyInput{
					SSHPublicKeyId: key.SSHPublicKeyId,
					UserName:       &user,
					Encoding:       "SSH",
				})
			if err != nil {
				return nil, err
			}
			keys = append(keys, *resp.SSHPublicKey.SSHPublicKeyBody)
		}
	}
	return keys, nil
}

// GetGroupUsers returns a list of users inside a group
func getGroupUsers(ctx context.Context, cfg appConfig.Config, c *iam.Client, g []string) ([]string, error) {
	var users []string
	for _, group := range g {
		// Get users from group
		resp, err := c.GetGroup(ctx, &iam.GetGroupInput{GroupName: &group})
		if err != nil {
			return nil, err
		}
		for _, user := range resp.Users {
			users = append(users, *user.UserName)
		}
	}
	return users, nil
}

// GetSSHKeys returns a list of SSH keys in total
func GetSSHKeys(c appConfig.Config) []string {
	ctx := context.TODO()

	log.Println("Loading AWS configuration using method " + c.Aws.Method)
	awsConfig := loadAwsConfig(ctx, c)
	client := getIAMClient(awsConfig)

	log.Println("Getting users from groups")
	users, err := getGroupUsers(ctx, c, client, c.Aws.Groups)
	if err != nil {
		log.Panicf("Error getting users from groups: %v\n", err)
	}
	log.Printf("Loading keys for users %v\n", len(users))
	log.Println("Getting SSH keys for users")
	keys, err := getKeysForUser(ctx, client, users)
	if err != nil {
		log.Panicf("Error getting SSH keys for users: %v\n", err)
	}
	log.Printf("Loaded %v keys", len(keys))
	ctx.Done()
	return keys
}
