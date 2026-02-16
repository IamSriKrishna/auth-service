#!/bin/bash

# DNS Configuration for MySQL Service
# This script creates DNS A record for MySQL service pointing to Traefik LoadBalancer

set -e

# Configuration
ZONE_NAME="bbcloud-app"
DOMAIN="bbcloud.app"
SUBDOMAIN="dev-mysql"
FULL_DOMAIN="${SUBDOMAIN}.${DOMAIN}"
TTL="300"

# Get Traefik LoadBalancer IP dynamically
echo "🔍 Getting Traefik LoadBalancer IP..."
TRAEFIK_IP=$(kubectl get service traefik -n traefik-system -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

if [ -z "$TRAEFIK_IP" ]; then
    echo "❌ Failed to get Traefik LoadBalancer IP"
    echo "🔧 Make sure Traefik is deployed and has an external IP"
    exit 1
fi

echo "🗄️ Setting up DNS for MySQL Service"
echo "Zone: ${ZONE_NAME}"
echo "Domain: ${FULL_DOMAIN}"
echo "IP: ${TRAEFIK_IP}"
echo "TTL: ${TTL} seconds"
echo ""

# Check if record already exists
echo "📋 Checking existing DNS records..."
if gcloud dns record-sets list --zone=${ZONE_NAME} --name=${FULL_DOMAIN}. --type=A --format="value(name)" | grep -q "${FULL_DOMAIN}."; then
    echo "⚠️  A record for ${FULL_DOMAIN} already exists. Updating..."
    
    # Get current IP to delete old record
    CURRENT_IP=$(gcloud dns record-sets list --zone=${ZONE_NAME} --name=${FULL_DOMAIN}. --type=A --format="value(rrdatas[0])")
    echo "Current IP: ${CURRENT_IP}"
    
    # Start transaction
    gcloud dns record-sets transaction start --zone=${ZONE_NAME}
    
    # Delete old record
    gcloud dns record-sets transaction remove --zone=${ZONE_NAME} \
        --name=${FULL_DOMAIN}. \
        --type=A \
        --ttl=${TTL} \
        ${CURRENT_IP}
    
    # Add new record
    gcloud dns record-sets transaction add --zone=${ZONE_NAME} \
        --name=${FULL_DOMAIN}. \
        --type=A \
        --ttl=${TTL} \
        ${TRAEFIK_IP}
    
    # Execute transaction
    gcloud dns record-sets transaction execute --zone=${ZONE_NAME}
    
    echo "✅ Updated A record for ${FULL_DOMAIN} -> ${TRAEFIK_IP}"
else
    echo "🆕 Creating new A record..."
    
    # Create new A record
    gcloud dns record-sets create ${FULL_DOMAIN}. \
        --zone=${ZONE_NAME} \
        --type=A \
        --ttl=${TTL} \
        --rrdatas=${TRAEFIK_IP}
    
    echo "✅ Created A record for ${FULL_DOMAIN} -> ${TRAEFIK_IP}"
fi

echo ""
echo "🔍 Verifying DNS record..."
gcloud dns record-sets list --zone=${ZONE_NAME} --name=${FULL_DOMAIN}. --type=A

echo ""
echo "🌍 DNS propagation may take a few minutes."
echo "🗄️ MySQL will be available at: ${FULL_DOMAIN}:3306"
echo "📝 Next: Create MySQL ingress to expose the service"
