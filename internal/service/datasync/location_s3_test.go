// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasync_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/datasync"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdatasync "github.com/hashicorp/terraform-provider-aws/internal/service/datasync"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDataSyncLocationS3_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationS3Output
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_datasync_location_s3.test"
	s3BucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationS3Destroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationS3Config_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "agent_arns.#", "0"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexp.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "s3_bucket_arn", s3BucketResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "s3_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_config.0.bucket_access_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttrSet(resourceName, "s3_storage_class"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/test/"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexp.MustCompile(`^s3://.+/`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_bucket_arn"},
			},
		},
	})
}

func TestAccDataSyncLocationS3_storageClass(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationS3Output
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	iamRoleResourceName := "aws_iam_role.test"
	resourceName := "aws_datasync_location_s3.test"
	s3BucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationS3Destroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationS3Config_storageClass(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexp.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "s3_bucket_arn", s3BucketResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "s3_config.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "s3_config.0.bucket_access_role_arn", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/test/"),
					resource.TestCheckResourceAttr(resourceName, "s3_storage_class", "STANDARD_IA"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexp.MustCompile(`^s3://.+/`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_bucket_arn"},
			},
		},
	})
}

func TestAccDataSyncLocationS3_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationS3Output
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_s3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationS3Destroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationS3Config_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdatasync.ResourceLocationS3(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSyncLocationS3_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationS3Output
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_s3.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationS3Destroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationS3Config_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"s3_bucket_arn"},
			},
			{
				Config: testAccLocationS3Config_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLocationS3Config_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationS3Exists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testAccCheckLocationS3Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_datasync_location_s3" {
				continue
			}

			_, err := tfdatasync.FindLocationS3ByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DataSync Location S3 %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckLocationS3Exists(ctx context.Context, n string, v *datasync.DescribeLocationS3Output) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		output, err := tfdatasync.FindLocationS3ByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccLocationS3Config_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "datasync.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "test" {
  role   = aws_iam_role.test.id
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [{
    "Action": [
      "s3:*"
    ],
    "Effect": "Allow",
    "Resource": [
      "${aws_s3_bucket.test.arn}",
      "${aws_s3_bucket.test.arn}/*"
    ]
  }]
}
POLICY
}

resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true
}
`, rName)
}

func testAccLocationS3Config_basic(rName string) string {
	return acctest.ConfigCompose(testAccLocationS3Config_base(rName), `
resource "aws_datasync_location_s3" "test" {
  s3_bucket_arn = aws_s3_bucket.test.arn
  subdirectory  = "/test"

  s3_config {
    bucket_access_role_arn = aws_iam_role.test.arn
  }
}
`)
}

func testAccLocationS3Config_storageClass(rName string) string {
	return acctest.ConfigCompose(testAccLocationS3Config_base(rName), `
resource "aws_datasync_location_s3" "test" {
  s3_bucket_arn    = aws_s3_bucket.test.arn
  subdirectory     = "/test"
  s3_storage_class = "STANDARD_IA"

  s3_config {
    bucket_access_role_arn = aws_iam_role.test.arn
  }
}
`)
}

func testAccLocationS3Config_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(testAccLocationS3Config_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_s3" "test" {
  s3_bucket_arn = aws_s3_bucket.test.arn
  subdirectory  = "/test"

  s3_config {
    bucket_access_role_arn = aws_iam_role.test.arn
  }

  tags = {
    %[1]q = %[2]q
  }
}
`, key1, value1))
}

func testAccLocationS3Config_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(testAccLocationS3Config_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_s3" "test" {
  s3_bucket_arn = aws_s3_bucket.test.arn
  subdirectory  = "/test"

  s3_config {
    bucket_access_role_arn = aws_iam_role.test.arn
  }

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, key1, value1, key2, value2))
}
