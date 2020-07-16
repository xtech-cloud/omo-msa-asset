package config

const defaultYAML string = `
service:
    address: :7079
    ttl: 15
    interval: 10
logger:
    level: info
    dir: /var/log/msa/
database:
    name: rgsCloud
    ip: 127.0.0.1
    port: "27017"
    user: root
    password: pass2019
    type: mongodb
basic:
    synonym: 6
    tag: 6
    
`
