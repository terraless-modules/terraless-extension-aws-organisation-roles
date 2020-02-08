package main

import "bytes"

const (
  accountAssumeRolePolicy = `
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
`

  accountRole = `
resource "aws_iam_role" "{{ .RoleName }}" {
  name = "{{ .RoleName }}"
  assume_role_policy = data.aws_iam_policy_document.assume-role.json
}

{{ range $idx, $val := .PolicyAttachments }}
resource "aws_iam_role_policy_attachment" "{{ $val.RoleName }}-{{ $val.Policy }}" {
  policy_arn = aws_iam_policy.{{ $val.Policy }}.arn
  role = aws_iam_role.{{ $val.RoleName }}.name
}

{{ end -}}
`

  accountPolicies = `
data "aws_iam_policy_document" "policy-{{ .PolicyName }}" {
{{ range $idx, $val := .Statements }}
  statement {
    effect = "{{ $val.Effect }}"

    actions = [
      {{ range $idx2, $val2 := $val.Actions }}"{{ $val2 }}",
      {{ end -}}
    ]

    resources = [
      {{ range $idx2, $val2 := $val.Resources }}"{{ $val2 }}",
      {{ end -}}
    ]
  }

{{ end -}}
}

resource "aws_iam_policy" "{{ .PolicyName }}" {
  name = "{{ .PolicyName }}"
  policy = data.aws_iam_policy_document.policy-{{ .PolicyName }}.json
}
`
)

func generateAccountRoles(config *OrganisationRolesConfig) bytes.Buffer {
  buffer := bytes.Buffer{}

  buffer.WriteString("# Generated by Terraless Organisation Roles - generateAccountRoles")
  buffer = renderToBuffer(config, accountAssumeRolePolicy, "assume-role-policy", buffer)

  for roleName, roleData := range config.Roles {
    var data = map[string]interface{} {
      "RoleName": roleName,
      "PolicyAttachments": policyData(roleName, roleData.Policies),
    }
    buffer = renderToBuffer(data, accountRole, "account-role", buffer)
  }

  for policyName, policyData := range config.Policies {
    var data = map[string]interface{} {
      "PolicyName": policyName,
      "Statements": policyData.Statements,
    }
    buffer = renderToBuffer(data, accountPolicies, "account-policies", buffer)
  }

  return buffer
}

func policyData(roleName string, policies []string) []map[string]string {
  var result []map[string]string

  for _, policy := range policies {
    data := map[string]string {
      "RoleName": roleName,
      "Policy": policy,
    }

    result = append(result, data)
  }

  return result
}