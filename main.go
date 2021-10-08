package main

import (
	"fmt"
	"os"

	"github.com/gofrs/uuid"
)

func main() {
	path := uuid.Must(uuid.NewV4()).String()
	controlPath := uuid.Must(uuid.NewV4()).String()
	os.WriteFile("/etc/nginx/sites-enabled/public", []byte(fmt.Sprintf(`server {
		listen 443 ssl;
		listen [::]:443 ssl;
	  
		access_log off;
		error_log off;
	  
		ssl_certificate /root/.crt;
		ssl_certificate_key /root/.key;
		ssl_session_timeout 1d;
		ssl_session_cache shared:MozSSL:10m;
		ssl_session_tickets off;
	  
		ssl_protocols TLSv1.2 TLSv1.3;
		ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
		ssl_prefer_server_ciphers off;
	  
		location /%s {
			if ($http_upgrade != "websocket") {
				return 404;
			}
			proxy_redirect off;
			proxy_pass http://127.0.0.1:17001;
			proxy_http_version 1.1;
			proxy_set_header Upgrade $http_upgrade;
			proxy_set_header Connection "upgrade";
			proxy_set_header Host $host;
			proxy_set_header X-Real-IP $remote_addr;
			proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		}
		location /%s/ {
			proxy_pass http://127.0.0.1:10087/;
		}
	}`, path, controlPath)), 0777)
	os.WriteFile("/usr/local/etc/v2ray/config.json", []byte(fmt.Sprintf(`{
		"stats": {},
		"api": {
			"tag": "api",
			"services": ["HandlerService", "LoggerService", "StatsService"]
		},
		"policy": {
			"levels": {
				"0": {
					"statsUserUplink": true,
					"statsUserDownlink": true
				}
			},
			"system": {
				"statsInboundUplink": true,
				"statsInboundDownlink": true,
				"statsOutboundUplink": true,
				"statsOutboundDownlink": true
			}
		},
		"inbounds": [
			{
				"tag": "proxy",
				"port": 17001,
				"listen":"127.0.0.1",
				"protocol": "vmess",
				"streamSettings": {
					"network": "ws",
					"wsSettings": {
						"path": "/%s"
					}
				}
			},
			{
				"tag": "api",
				"listen": "127.0.0.1",
				"port": 10085,
				"protocol": "dokodemo-door",
				"settings": {
					"address": "127.0.0.1"
				}
			}
		],
		"outbounds": [
			{
				"protocol": "freedom",
				"settings": {}
			},
			{
				"tag": "block",
				"protocol": "blackhole",
				"settings": {}
			}
		],
		"routing": {
			"rules": [
				{
					"type": "field",
					"inboundTag": ["api"],
					"outboundTag": "api"
				},
				{
					"type": "field",
					"protocol": ["bittorrent"],
					"outboundTag": "block"
				}
			]
		}
	}`, path)), 0777)
	fmt.Println(path, controlPath)
}
