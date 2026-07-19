# imap-delete

A CLI tool to automatically delete messages from an IMAP mailbox.

## Installation

```sh
curl -Lo imap-delete https://github.com/analike/go-imap-delete/releases/download/v1.0.0/imap-delete-v1.0.0-linux-amd64 && \
chmod +x imap-delete
```

## Configuration

Create an `imap.yaml` file using the `imap.example.yaml` template:

```yaml
timeout:
  connect: 10 # seconds
  command: 30 # seconds
  delay: 300 # milliseconds to arbitrarily wait between iterations for monitoring
mailboxes:
  project-1:
    host: mx.project-1.com
    port: 993
    secure: true
    auth:
      user: me@project-1.com
      pass: 'xxx111222333'
  personal:
    host: imap.personal.com
    port: 993
    secure: true
    auth:
      user: me@personal.com
      pass: 'yyy22233334444'
  work:
    host: imap.work.com
    port: 993
    secure: true
    auth:
      user: team@work.com
      pass: 'zzz33344445555'
```

Using this example, `project-1`, `personal` and `work` are the mailbox names for running the tool.

## Usage

```sh
./imap-delete project-1 --folder 'Junk' --subject 'status report'
```

### Options

The positional `[box]` (mailbox name as defined in config yaml) argument is required, all others are optional.

::: info
Dry run is **off** by default. It is only auto-enabled when no filters are specified and `--force` is not passed. This is a safety net to prevent accidentally wiping a whole folder. Pass `--force` to bypass this and run for real with no filters, or pass at least one filter option to run normally.
:::

When multiple filters are given, they are combined with **AND** logic — a message must match all specified filters to be deleted. For example:

```sh
./imap-delete project-1 --folder 'Junk' --subject 'status report' --before '2025-01-01'
```

This only deletes messages that are in `Junk`, **and** match the subject, **and** were received before `2025-01-01`.

All output and errors (e.g. auth failures, connection timeouts) are logged to the terminal - useful for wrapping this in cron or CI.

#### 1. Config Options

|        Option | Description                                                                                              | Default          |
|--------------:|----------------------------------------------------------------------------------------------------------|------------------|
|   **\[box\]** | First positional argument specifies the config to select from the list of connections in the config file | **\*\*REQUIRED** |
|  **--config** | Path to config file                                                                                      | `./imap.yaml`    |
|  **--folder** | Which folder to search for messages \[ e.g `--folder 'Junk'` \]                                          | `INBOX`          |
| **--dry-run** | Boolean flag. Search and print messages only, without deleting \[ e.g `--dry-run` \]                     | `false`          |
|   **--force** | Used to bypass the safety mode of not deleting messages when no filters are specified                    | `false`          |

#### 2. Filter Options

Multiple filters are combined with **AND** logic.

|         Option | Description                                                                                                                  | Default |
|---------------:|------------------------------------------------------------------------------------------------------------------------------|---------|
|      **--uid** | Specify message `uid` to search                                                                                              | `NULL`  |
|       **--to** | Search in the recipient address field \[e.g `--to 'Peter Gonzalez'` \]                                                       | `NULL`  |
|       **--cc** | Search in the `cc` field \[e.g `--cc 'John Lambert'` \]                                                                      | `NULL`  |
|      **--bcc** | Search in the `bcc` field \[e.g `--bcc 'lina.rosario@gmail.com'` \]                                                          | `NULL`  |
|     **--from** | Search in the `from` header \[e.g `--from 'Peter Smith'` \]                                                                  | `NULL`  |
|  **--subject** | Search in the subject field \[e.g `--subject 'Transaction Notification'` \]                                                  | `NULL`  |
|     **--body** | Search in message body \[e.g `--body 'new sign in to your'` \]                                                               | `NULL`  |
|     **--text** | Search in message header or body \[e.g `--text 'commented on a'` \]                                                          | `NULL`  |
|    **--since** | Search for messages received on or after `value` \[ `--since '2025-01-01'` \]                                                | `NULL`  |
|   **--before** | Search for messages received before `value` \[ `--before '2025-12-12'` \]                                                    | `NULL`  |
|     **--date** | Search for messages received on a specific date \[ `--date '2025-06-30'` \]                                                  | `NULL`  |
| **--answered** | Search for messages that have been answered (i.e `\Answered` flag). Accepts `yes`/`no`/`true`/`false` \[ `--answered yes` \] | `NULL`  |
|  **--deleted** | Search for messages that are marked as deleted (i.e `\Deleted` flag). Accepts `yes`/`no`/`true`/`false` \[ `--deleted no` \] | `NULL`  |
|    **--draft** | Search for messages that are marked as draft (i.e `\Draft` flag). Accepts `yes`/`no`/`true`/`false` \[ `--draft yes` \]      | `NULL`  |
|  **--flagged** | Search for messages that are marked as flagged (i.e `\Flagged` flag). Accepts `yes`/`no`/`true`/`false` \[ `--flagged no` \] | `NULL`  |
|     **--seen** | Search for messages that are marked as seen (i.e `\Seen` flag). Accepts `yes`/`no`/`true`/`false` \[ `--seen no` \]          | `NULL`  |

## License

MIT