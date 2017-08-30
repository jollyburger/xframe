{{.AppName}}:
  working_dir: /opt
  image: 
  container_name: {{.AppName}}
  command: 
  net: "host"
  ports:
  - "9096:9096"
  volumes: 
  - "/opt/{{.AppName}}/config.json:/opt/config.json"