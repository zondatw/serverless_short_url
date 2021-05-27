package shorturlsourcefunction

var AVRO_SOURCE string = `{
    "type": "record",
    "name": "Avro",
    "fields": [
      {
        "name": "Datetime",
        "type": "string"
      },
      {
        "name": "SourceIp",
        "type": "string"
      },
      {
        "name": "Agent",
        "type": "string"
      },
      {
        "name": "ShortHash",
        "type": "string"
      }
    ]
  }`
