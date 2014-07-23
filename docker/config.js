var api_names = [];
var apis = [];
for (var prop in process.env) { if (prop.match(/PORT_4567_TCP_ADDR/)) { var res = prop.split('_'); api_names.push(res[0]) } };
for (var n in api_names) {
  var name = api_names[n]
  var api = {}
  api['name'] = process.env[name+"_UCHIWA_NAME"] || name
  api['host'] = process.env[name+"_PORT_4567_TCP_ADDR"]
  api['port'] = parseInt(process.env[name+"_PORT_4567_TCP_PORT"])
  api['ssl'] = process.env[name+"_SSL"] ? true : false
  api['user'] = process.env[name+"_UCHIWA_USER"] || ''
  api['pass'] = process.env[name+"_UCHIWA_PASS"] || ''
  api['path'] = process.env[name+"_UCHIWA_PATH"] || ''
  api['timeout'] = parseInt(process.env[name+"_UCHIWA_TIMEOUT"]) || 5000
  apis.push(api)
}
var allconf = {
  sensu: apis,
  uchiwa: {
    host: '0.0.0.0',
    port: 3000,
    user: process.env.UCHIWA_USER || '',
    pass: process.env.UCHIWA_PASS || '',
    stats: parseInt(process.env.UCHIWA_STATS),
    refresh: parseInt(process.env.UCHIWA_REFRESH)
  }
}
console.log(allconf)
module.exports = allconf
