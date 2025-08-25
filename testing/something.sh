# Run the terraform state list command and filter resources by "newrelic_nrql_drop_rule"
resources=$(terraform state list | grep "newrelic_nrql_drop_rule")


json_output="{\"resources\":\""

# Loop through each resource and add it to the JSON string
for resource in $resources; do
  resource_identifier=$(echo "$resource" | sed 's/.*newrelic_nrql_drop_rule\.//')
  json_output+="$resource_identifier,"
done

# Remove the trailing comma and close the JSON string
json_output=$(echo "$json_output" | sed 's/,$//')
json_output+="\"}"

# Print the JSON object
echo "$json_output"