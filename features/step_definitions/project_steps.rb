# Givens

Given(/^I have a project$/) do
  steps %Q{
    Given I have a ruleset id
    When I run the command to create a project
  }

  json = JSON.parse(@output)
  @project_id = json['id']
end

# Whens

When(/^I run the command to create a project$/) do
  team = 'test-team'
  project_name = 'java-lew'
  @branch = "master-#{Time.now.to_i.to_s + rand(1000).to_s}"
  source = 'git@github.com:ion-channel/java-lew.git'
  description = 'Java Lew'

  @output = `ion-connect project create-project --team-id #{team} --ruleset-id #{@ruleset_id} --active --branch #{@branch} #{project_name} "#{source}" "#{description}"`.chomp
end

When(/^I run the command to get the project$/) do
  team = 'test-team'

  @output = `ion-connect project get-project --team-id #{team} #{@project_id}`.chomp
end

# Thens

Then(/^I see a response showing the project is created$/) do
  expect(@output).to include('"active": true')
  expect(@output).to include('"branch": "master-')
  expect(@output).to include('"source": "git@github.com:ion-channel/java-lew.git"')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end

Then(/^I see the projects details$/) do
  expect(@output).to include('"active": true')
  expect(@output).to include('"branch": "master-')
  expect(@output).to include('"source": "git@github.com:ion-channel/java-lew.git"')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end
