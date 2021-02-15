package internal

import (
    "fmt"
    // "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sts"
)

type Credentials struct {
    sts sts.Credentials
    profile string
}

// SourceableBashEnv prints commands to set credentials
// can be used in combination with $eval() to affect parent shell
func SourceableBashEnv(c *Credentials) {
    fmt.Printf("set -x DD_ASSUMED_ROLE %s\n", c.profile)
    fmt.Printf("set -x AWS_ACCESS_KEY_ID %s\n", *c.sts.AccessKeyId)
    fmt.Printf("set -x AWS_SECRET_ACCESS_KEY %s\n", *c.sts.SecretAccessKey)
    fmt.Printf("set -x AWS_SESSION_TOKEN %s\n", *c.sts.SessionToken)
    if (c.sts.Expiration != nil) {
        fmt.Printf("set -lx AWS_EXPIRATION %s\n", c.sts.Expiration)
    }
}

// SourceableUnsetBashEnv prints commands to unset credentials
// can be used in combination with $eval() to affect parent shell
func SourceableUnsetBashEnv() {
    fmt.Println("set -e DD_ASSUMED_ROLE")
    fmt.Println("set -e AWS_ACCESS_KEY_ID\n")
    fmt.Println("set -e AWS_SECRET_ACCESS_KEY\n")
    fmt.Println("set -e AWS_SESSION_TOKEN\n")
}

// AssumeRoleViaProfile attempts to acquire temporary role credentials using AWS config settings paired with a profile name
// respects default chain of credential providers - i.e. env, shared credentials file (~/.aws/credentials) or EC2 instance role
func AssumeRoleViaProfile(profile string) (*Credentials, error) {

    // create session
    sess, err := session.NewSessionWithOptions(session.Options{
        Profile:           profile,
        SharedConfigState: session.SharedConfigEnable,
        // AssumeRoleDuration: time.Duration(900),
    })

    if err != nil {
        return nil, err
    }

    result, err := sess.Config.Credentials.Get()

    // force result into (currently) from credentials.Value into standardized upon sts.Credentials
    var c sts.Credentials
    c.AccessKeyId = aws.String(result.AccessKeyID)
    c.SecretAccessKey = aws.String(result.SecretAccessKey)
    c.SessionToken = aws.String(result.SessionToken)

    var cc Credentials
    cc.sts = c
    cc.profile = profile

    return &cc, nil
}
