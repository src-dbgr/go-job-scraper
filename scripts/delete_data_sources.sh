#!/bin/bash

# Enable debug mode
set -x

# Prompt for Grafana API key
read -p "Enter your Grafana API key: " API_KEY

# Prompt for Grafana URL with default value
read -p "Enter Grafana URL [http://localhost:3000]: " GRAFANA_URL
GRAFANA_URL=${GRAFANA_URL:-http://localhost:3000}

# Test API connection and get all datasources
echo "Testing API connection and fetching all datasources..."
all_datasources=$(curl -s -H "Authorization: Bearer $API_KEY" "${GRAFANA_URL}/api/datasources")
echo "All datasources:"
echo "$all_datasources" | jq '.'

# Function to delete a data source
delete_datasource() {
    local name=$1
    local full_name="v1 ${name}"
    
    echo "Checking if datasource exists: ${full_name}"
    # Find the datasource with the exact name
    datasource=$(echo "$all_datasources" | jq -r ".[] | select(.name == \"$full_name\")")
    
    if [ -n "$datasource" ]; then
        id=$(echo "$datasource" | jq -r '.id')
        echo "Found datasource: ${full_name} with ID: ${id}"
        
        # Delete existing datasource
        delete_response=$(curl -s -X DELETE -H "Authorization: Bearer $API_KEY" \
        "${GRAFANA_URL}/api/datasources/${id}")
        echo "Delete response for ${full_name}:"
        echo "$delete_response" | jq '.'
    else
        echo "Data source not found: ${full_name}"
    fi
}

# Main execution
echo "Choose an action:"
echo "1. Create all v1 data sources"
echo "2. Delete all v1 data sources"
read -p "Enter your choice (1 or 2): " choice

case $choice in
    2)
        # List of data source names (without the "v1 " prefix)
        datasources=(
            "avg salary by education"
            "avg-experience-by-category"
            "benefits by company size"
            "companies-by-size"
            "company size distribution"
            "employment types"
            "Job category count"
            "job postings per company"
            "job postings per day"
            "job postings per month"
            "job postings trend"
            "job-categories-by-company-size"
            "job-requirements-by-location"
            "remote vs onsite"
            "remote vs onsite by industry"
            "skills-by-experience-level"
            "tech trends"
            "top job category"
            "top optional skills"
            "top skills"
        )

        for ds in "${datasources[@]}"; do
            delete_datasource "$ds"
        done
        echo "All v1 data sources have been processed."
        ;;
    *)
        echo "Invalid choice or not implemented. Exiting."
        exit 1
        ;;
esac