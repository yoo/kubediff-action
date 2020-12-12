# KubeDiff Action

This action runs `kubectl diff` for a pull request and adds the as a comment to the pull request.
To see KubeDiff in action have look at the [Demo Pull Request](https://github.com/yoo/kubediff-action/pull/1).

## Example

```yaml
- name: KubeDiff
  uses: yoo/kubediff-action@v1
  with:
    manifests: |
      "manifest1.yml"
      "manifest2.yml"
```

## Inputs

```yaml
inputs:
  manifests:
    description: "Line separated list of manifest files"
    required: true
  comment_pr:
    description: "Create a pull request comment containing the diff"
    default: "true"
  filtered_fields:
    description: "Line separated list of fields to filter out when displaying the diff"
    required: true
    default: |
      /metadata/annotations/kubectl.kubernetes.io\/last-applied-configuration
      /metadata/creationTimestamp
      /metadata/managedFields
      /metadata/resourceVersion
      /metadata/selfLink
      /metadata/uid
  debug:
    description: "Toggle debug output"
    required: false
```

