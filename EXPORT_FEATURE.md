# Pull Request Export Feature

This feature allows you to export pull requests from a GitHub repository to a text file using GraphQL.

## Usage

Use the `export` field on the `GitHubPullRequestType` with a `repositoryURL` argument:

```graphql
query {
  github {
    pullRequests {
      export(repositoryURL: "https://github.com/owner/repo-name")
    }
  }
}
```

## Output Format

The export creates a file at `assets/-<repo-name>/exports/pull_requests/data.txt` with pull request data in pipe-separated format:

```
Number|Title|Contributors|Duration|CreatedAt|MergedAt|AbbreviatedBody
```

### Example Output

```
123|Fix authentication bug|user1, user2|52h30m0s|2024-01-01T10:00:00Z|2024-01-03T14:30:00Z|This pull request fixes a critical authentication bug that was causing users to be unable to log in. The issue was in the JWT token validation logic w ...
124|Add new feature for user profiles|developer1|55h45m0s|2024-01-05T09:00:00Z|2024-01-07T16:45:00Z|Implementing user profile functionality with avatar upload support.
```

## Field Descriptions

- **Number**: Pull request number
- **Title**: Pull request title
- **Contributors**: Comma-separated list of contributor usernames
- **Duration**: Time between creation and merge (MergedAt - CreatedAt)
- **CreatedAt**: Pull request creation time in RFC3339 format
- **MergedAt**: Pull request merge time in RFC3339 format
- **AbbreviatedBody**: First 150 characters of the pull request body

## File Handling

- If the export file already exists, the operation returns a message without overwriting
- Only merged pull requests (with both CreatedAt and MergedAt) are included
- The directory structure is automatically created if it doesn't exist