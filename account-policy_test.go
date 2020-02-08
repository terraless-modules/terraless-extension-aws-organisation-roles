package main

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestAccountPolicy_generateAccountRoles(t *testing.T) {
  // given
  config := OrganisationRolesConfig{
    Variables: Variables{
      MainAccountId: "235324543534",
    },
    Roles: map[string]Role{
      "admin": {
        Policies: []string{
          "admin",
          "developer",
          "billing-readonly",
        },
      },
      "developer": {
        Policies: []string{
          "developer",
          "billing-readonly",
        },
      },
    },
    Policies: map[string]Policy{
      "admin": {
        Statements: []Statement{
          {
            Effect: "Allow",
            Resources: []string{
              "resource1",
              "resource2",
            },
            Actions: []string{
              "action1",
              "action2",
            },
          },
        },
      },
      "developer": {
        Statements: []Statement{
          {
            Effect: "Allow",
            Resources: []string{
              "resource1",
              "resource2",
            },
            Actions: []string{
              "action1",
              "action2",
            },
          },
          {
            Effect: "Deny",
            Resources: []string{
              "resource1",
              "resource3",
            },
            Actions: []string{
              "action4",
            },
          },
        },
      },
    },
  }

  // when
  result := generateAccountRoles(&config)

  // then
  expectedResult := `# Generated by Terraless Organisation Roles - generateAccountRoles
data "aws_iam_policy_document" "assume-role" {
  statement {
    effect = "Allow"

    actions = ["sts:AssumeRole"]

    principals {
      identifiers = [
        "arn:aws:iam::${var.root-account-id}:root",
      ]
      type = "AWS"
    }

    condition {
      test = "Bool"
      values = ["true"]
      variable = "aws:MultiFactorAuthPresent"
    }
  }
}

variable "root-account-id" {
	type = string
}

resource "aws_iam_role" "admin" {
  name = "admin"
  assume_role_policy = data.aws_iam_policy_document.assume-role.json
}


resource "aws_iam_role_policy_attachment" "admin-admin" {
  policy_arn = aws_iam_policy.admin.arn
  role = aws_iam_role.admin.name
}


resource "aws_iam_role_policy_attachment" "admin-developer" {
  policy_arn = aws_iam_policy.developer.arn
  role = aws_iam_role.admin.name
}


resource "aws_iam_role_policy_attachment" "admin-billing-readonly" {
  policy_arn = aws_iam_policy.billing-readonly.arn
  role = aws_iam_role.admin.name
}


resource "aws_iam_role" "developer" {
  name = "developer"
  assume_role_policy = data.aws_iam_policy_document.assume-role.json
}


resource "aws_iam_role_policy_attachment" "developer-developer" {
  policy_arn = aws_iam_policy.developer.arn
  role = aws_iam_role.developer.name
}


resource "aws_iam_role_policy_attachment" "developer-billing-readonly" {
  policy_arn = aws_iam_policy.billing-readonly.arn
  role = aws_iam_role.developer.name
}


data "aws_iam_policy_document" "policy-admin" {

  statement {
    effect = "Allow"

    actions = [
      "action1",
      "action2",
      ]

    resources = [
      "resource1",
      "resource2",
      ]
  }

}

resource "aws_iam_policy" "admin" {
  name = "admin"
  policy = data.aws_iam_policy_document.policy-admin.json
}

data "aws_iam_policy_document" "policy-developer" {

  statement {
    effect = "Allow"

    actions = [
      "action1",
      "action2",
      ]

    resources = [
      "resource1",
      "resource2",
      ]
  }


  statement {
    effect = "Deny"

    actions = [
      "action4",
      ]

    resources = [
      "resource1",
      "resource3",
      ]
  }

}

resource "aws_iam_policy" "developer" {
  name = "developer"
  policy = data.aws_iam_policy_document.policy-developer.json
}
`

  assert.Equal(t, expectedResult, result.String())
}