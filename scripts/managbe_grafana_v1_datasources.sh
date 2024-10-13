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

# Function to create a data source
create_datasource() {
    local name=$1
    local endpoint=$2
    local full_name="v1 ${name}"
    
    echo "Checking if datasource exists: ${full_name}"
    # Find the datasource with the exact name
    existing_datasource=$(echo "$all_datasources" | jq -r ".[] | select(.name == \"$full_name\")")
    
    if [ -n "$existing_datasource" ]; then
        echo "Data source already exists: ${full_name}. Skipping."
    else
        echo "Creating new datasource: ${full_name}"
        create_response=$(curl -s -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $API_KEY" \
        "${GRAFANA_URL}/api/datasources" \
        -d '{
          "name": "'"${full_name}"'",
          "type": "marcusolsson-json-datasource",
          "url": "http://host.docker.internal:8080/api/v1/'"${endpoint}"'",
          "access": "proxy",
          "basicAuth": false
        }')
        echo "Create response for ${full_name}:"
        echo "$create_response" | jq '.'
    fi
}

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

# List of datasource names and endpoints
datasources="
avg salary by education:stats/avg-salary-by-education
avg-experience-by-category:stats/avg-experience-by-category
benefits by company size:stats/benefits-by-company-size
companies-by-size:stats/companies-by-size
company size distribution:stats/company-size-distribution
employment types:stats/employment-types
Job category count:stats/job-categories-counts
job postings per company:stats/job-postings-per-company
job postings per day:stats/job-postings-per-day
job postings per month:stats/job-postings-per-month
job postings trend:stats/job-postings-trend
job-categories-by-company-size:stats/job-categories-by-company-size
job-requirements-by-location:stats/job-requirements-by-location
remote vs onsite:stats/remote-vs-onsite
remote vs onsite by industry:stats/remote-vs-onsite-by-industry
skills-by-experience-level:stats/skills-by-experience-level
tech trends:stats/technology-trends
top job category:stats/top-job-categories
top optional skills:stats/top-optional-skills
top skills:stats/top-skills
must skills:stats/mustskills
optional skills:stats/optionalskills
"

case $choice in
    1)
        echo "$datasources" | while IFS=':' read -r name endpoint; do
            create_datasource "$name" "$endpoint"
        done
        echo "All v1 data sources have been created or updated."
        ;;
    2)
        echo "$datasources" | while IFS=':' read -r name endpoint; do
            delete_datasource "$name"
        done
        echo "All v1 data sources have been processed for deletion."
        ;;
    *)
        echo "Invalid choice. Exiting."
        exit 1
        ;;
esac