// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package kinesisvideo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisvideo"
	"github.com/aws/aws-sdk-go/service/kinesisvideo/kinesisvideoiface"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists kinesisvideo service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(conn kinesisvideoiface.KinesisVideoAPI, identifier string) (tftags.KeyValueTags, error) {
	return ListTagsWithContext(context.Background(), conn, identifier)
}

func ListTagsWithContext(ctx context.Context, conn kinesisvideoiface.KinesisVideoAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &kinesisvideo.ListTagsForStreamInput{
		StreamARN: aws.String(identifier),
	}

	output, err := conn.ListTagsForStreamWithContext(ctx, input)

	if err != nil {
		return tftags.New(nil), err
	}

	return KeyValueTags(output.Tags), nil
}

// map[string]*string handling

// Tags returns kinesisvideo service tags.
func Tags(tags tftags.KeyValueTags) map[string]*string {
	return aws.StringMap(tags.Map())
}

// KeyValueTags creates KeyValueTags from kinesisvideo service tags.
func KeyValueTags(tags map[string]*string) tftags.KeyValueTags {
	return tftags.New(tags)
}

// UpdateTags updates kinesisvideo service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func UpdateTags(conn kinesisvideoiface.KinesisVideoAPI, identifier string, oldTags interface{}, newTags interface{}) error {
	return UpdateTagsWithContext(context.Background(), conn, identifier, oldTags, newTags)
}
func UpdateTagsWithContext(ctx context.Context, conn kinesisvideoiface.KinesisVideoAPI, identifier string, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := tftags.New(oldTagsMap)
	newTags := tftags.New(newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &kinesisvideo.UntagStreamInput{
			StreamARN:  aws.String(identifier),
			TagKeyList: aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagStreamWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &kinesisvideo.TagStreamInput{
			StreamARN: aws.String(identifier),
			Tags:      Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagStreamWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
