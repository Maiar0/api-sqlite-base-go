#!/usr/bin/env bash
set -euo pipefail

# ====== CONFIGURABLE SETTINGS ======
APP_NAME="api-sqlite-base-go"           # Go binary name
APP_USER="ubuntu"                       # Linux user that will run the service
APP_DIR="/opt/$APP_NAME"                # Install directory
SERVICE_NAME="$APP_NAME.service"
APP_PORT="3000"                         # Port your Go app listens on
DOMAIN_NAME="dwarddevs.com"             # Replace with your real domain if you have one
# ===================================

echo "==> Checking for root privileges..."
if [[ "$EUID" -ne 0 ]]; then
  echo "Please run this script as root, e.g.:"
  echo "  sudo ./deploy.sh"
  exit 1
fi

echo "==> Updating apt and installing packages (nginx, sqlite3, curl)..."
apt-get update -y
apt-get install -y nginx sqlite3 curl

echo "==> Creating application directory: $APP_DIR"
mkdir -p "$APP_DIR/store/auth"

echo "==> Copying binary into place..."
if [[ ! -f "./$APP_NAME" ]]; then
  echo "ERROR: Binary ./$APP_NAME not found in current directory."
  echo "Build it locally and upload it next to deploy.sh, then rerun."
  exit 1
fi

cp "./$APP_NAME" "$APP_DIR/$APP_NAME"
chmod +x "$APP_DIR/$APP_NAME"

echo "==> Setting ownership to $APP_USER..."
chown -R "$APP_USER":"$APP_USER" "$APP_DIR"

echo "==> Creating systemd service: /etc/systemd/system/$SERVICE_NAME"
cat >/etc/systemd/system/$SERVICE_NAME <<EOF
[Unit]
Description=Go API Service ($APP_NAME)
After=network.target

[Service]
Type=simple
User=$APP_USER
WorkingDirectory=$APP_DIR
ExecStart=$APP_DIR/$APP_NAME
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

echo "==> Reloading systemd and enabling service..."
systemctl daemon-reload
systemctl enable "$SERVICE_NAME"
systemctl restart "$SERVICE_NAME"

echo "==> Checking service status (short)..."
systemctl --no-pager --full status "$SERVICE_NAME" || true

echo "==> Configuring Nginx reverse proxy..."
NGINX_SITE="/etc/nginx/sites-available/api-sqlite-base-go"
cat >"$NGINX_SITE" <<EOF
server {
    listen 80;
    server_name $DOMAIN_NAME;

    location / {
        proxy_pass         http://127.0.0.1:$APP_PORT;
        proxy_http_version 1.1;
        proxy_set_header   Upgrade \$http_upgrade;
        proxy_set_header   Connection "upgrade";
        proxy_set_header   Host \$host;
        proxy_set_header   X-Real-IP \$remote_addr;
        proxy_set_header   X-Forwarded-For \$proxy_add_x_forwarded_for;
    }
}
EOF

ln -sf "$NGINX_SITE" /etc/nginx/sites-enabled/api-sqlite-base-go

echo "==> Testing Nginx configuration..."
nginx -t

echo "==> Restarting Nginx..."
systemctl restart nginx

echo "==> Deployment complete."

echo
echo "Summary:"
echo "  - Binary installed in: $APP_DIR"
echo "  - SQLite DB will be created in: $APP_DIR/auth/users.db"
echo "  - Systemd service: $SERVICE_NAME (user: $APP_USER)"
echo "  - Nginx is proxying http://YOUR_SERVER_IP/ -> 127.0.0.1:$APP_PORT"
echo
echo "Check logs with:"
echo "  journalctl -u $SERVICE_NAME -f"
