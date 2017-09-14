# Givens

Given(/^I have a ruleset id$/) do
  steps %Q{
    Given an Ion Channel team id 'team-id'
    When I successfully run with 'team_id' `ion-connect ruleset create-ruleset --team-id team_id test-ruleset "this is a test ruleset" '["c30b917956c3040daa2c571ef31dbe3a"]'`
  }

  json = JSON.parse($output)
  @ruleset_id = json['id']
end

# Whens

When(/^I run the command to create a project$/) do
  team = 'test-team'
  project_name = "java-lew-#{Time.now.to_i.to_s + rand(1000).to_s}"
  source = 'git@github.com:ion-channel/java-lew.git'
  description = 'Java Lew'

  @output = `ion-connect project create-project --team-id #{team} --ruleset-id #{@ruleset_id} --active #{project_name} "#{source}" "#{description}"`.chomp
end

Then(/^I see a response showing the project is created$/) do
  expect(@output).to include('"active": true')
  expect(@output).to include('"branch": "master"')
  expect(@output).to include('"source": "git@github.com:ion-channel/java-lew.git"')
end
