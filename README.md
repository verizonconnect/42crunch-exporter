# 42crunch-exporter

## Usage

```
usage: 42crunch-exporter --42c-api-key=$42C-API-KEY [<flags>]

Flags:
  -h, --[no-]help                Show context-sensitive help (also try
                                 --help-long and --help-man).
      --web.listen-address=:9916 ...
                                 Addresses on which to expose metrics and web
                                 interface. Repeatable for multiple addresses.
      --web.config.file=""       [EXPERIMENTAL] Path to configuration file
                                 that can enable TLS or authentication. See:
                                 https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
      --web.metrics-path="/metrics"
                                 Path under which to expose metrics
      --42c-address="https://platform.42crunch.com"
                                 42Crunch server address (can also be set with
                                 $42C_ADDR) ($42C_ADDR)
      --42c-api-key=42C-API-KEY  42Crunch API key (can also be set with
                                 $42C_API_KEY) ($42C_API_KEY)
      --42c-collection-regex=42C-COLLECTION-REGEX
                                 Regex which will include only specific
                                 42Crunhc API collections. (can also be set with
                                 $42C_COLLECTION_REGEX) ($42C_COLLECTION_REGEX)
      --log.level=info           Only log messages with the given severity or
                                 above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt,
                                 json]
      --[no-]version             Show application version.
```

## Metrics

| Metric                                          | Meaning                                                               | Labels                                           |
| ----------------------------------------------- | --------------------------------------------------------------------- | ------------------------------------------------ |
| fortytwo_crunch_collection_information            | Basic information about the api collection                             | id,name                      |
| fortytwo_crunch_api_information                   | Basic information about an API                                         | collection_id,id,name,tags   |
| fortytwo_crunch_api_assessment_criticals          | The number of critical vulnerabilities per api based on the API Audit  | id                           |
| fortytwo_crunch_api_assessment_highs              | The number of high vulnerabilities per api based on the API Audit      | id                           |
| fortytwo_crunch_api_assessment_mediums            | The number of medium vulnerabilities per api based on the API Audit    | id                           |
| fortytwo_crunch_api_assessment_lows               | The number of low vulnerabilities per api based on the API Audit       | id                           |
| fortytwo_crunch_api_assessment_infos              | The number of information messages per api based on the API Audit      | id                           |
| fortytwo_crunch_api_assessment_grade              | API Audit Assessment Grade                                             | id                           |
| fortytwo_crunch_api_assessment_errors             | The number of API errors                                               | id                           |
| fortytwo_crunch_api_assessment_valid              | Indicating whether the api schema is valid                             | id                           |
| fortytwo_crunch_api_assessment_last_audit         | Last API Audit Assessment date, represented as a Unix timestamp        | id                           |
| fortytwo_crunch_api_assessment_semantic_invalid   | Indicates whether the api has semantic issues                          | id                           |
| fortytwo_crunch_api_assessment_structure_invalid  | Indicates whether the api has structural issues                        | id                           |
