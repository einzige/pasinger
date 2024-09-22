require 'json'
require 'csv'

# Load the CSV file with eva_number;station_name
eva_csv_path = 'evas.csv'
eva_data = {}

CSV.foreach(eva_csv_path, col_sep: ';', headers: false) do |row|
  eva_number, station_name = row
  eva_data[station_name] = eva_number
end

# Load the JSON file
trainmap_json_path = 'trainmap_cells_corrected.json'
trainmap_data = JSON.parse(File.read(trainmap_json_path))

# Add EVA numbers to destination cells
trainmap_data.each do |cell|
  if cell['type'] == 'destination'
    station_name = cell['text']
    cell['eva'] = eva_data[station_name] || 'fixme'
  end
end

# Save the updated JSON data to a new file
output_json_path = 'trainmap_with_eva.json'
File.open(output_json_path, 'w') do |f|
  f.write(JSON.pretty_generate(trainmap_data))
end

puts "Updated JSON with EVA numbers saved to #{output_json_path}"
