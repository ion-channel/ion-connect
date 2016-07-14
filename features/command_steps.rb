When(/^I successfully run debug `(.*)`$/) do |command|
  $output = `#{command}`
  expect($?.to_i).to eql(0)
end

When(/^I successfully run ion `(.*)`$/) do |command|
  directory = Dir.pwd
  Dir.chdir 'tmp/aruba'
  @output = `#{command}`.chomp
  Dir.chdir directory
  expect($?.to_i).to eql(0)
end

When(/^I successfully run with '(.*)' `(.*)`$/) do |variables, command|
  variables.split(',').each do |variable|
    command = command.gsub(/\s#{variable}\s*/, " #{instance_variable_get("@#{variable}")} ")
  end

  Aruba.configure{ |config| dir = config.working_directory}
  @output = `#{command}`.chomp
  $output = @output

  expect($?.to_i).to eql(0)
end

Then(/^the ion output should contain:$/) do |string|
  expect(@output).to include(string)
end

Given(/^an Ion Channel account id '(.*)'/) do |account_id|
  @account_id = account_id
end

Given(/^a variable '(.*)' is set from the previous output from location '(.*)'$/) do |variable, location|
  json = JSON.parse($output)
  instance_variable_set("@#{variable}", json[location])
end

Given(/^previous output/) do
  expect($output).to_not be_nil
end
