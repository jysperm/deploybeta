local http = require 'resty.http'
local cjson = require 'cjson.safe'

local etcd = {}

function etcd.getBackends(domain)
  local httpc = http.new()

  local res, err = httpc:request_uri('http://127.0.0.1:2379/v3alpha/kv/range', {
    method = 'POST',
    headers = {
      ["Content-Type"]= 'application/json',
    },
    body = cjson.encode({
      key = ngx.encode_base64('/upstreams/' .. domain)
    })
  })

  local body = cjson.decode(res.body)

  if body['kvs'] && body['kvs'][1] then
    local backends = {}

    for index, key in ipairs(cjson.decode(ngx.decode_base64(body['kvs'][1]['value']))['backends']) do
      backends[index] = '127.0.0.1:' .. key.port
    end

    return backends
  else
    return nil
  end
end

return etcd
