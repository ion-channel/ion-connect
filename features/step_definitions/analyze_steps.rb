# Givens

Given(/^I have a project$/) do
  steps %Q{
    Given I have a ruleset id
    When I run the command to create a project
  }

  json = JSON.parse(@output)
  @project_id = json['id']
end

Given(/^I have a ruleset id$/) do
  steps %Q{
    Given I have a set of rules
    When I run the command to create a ruleset
  }

  json = JSON.parse(@output)
  @ruleset_id = json['id']
end

Given(/^I have a set of rules$/) do
  @rules = ['c30b917956c3040daa2c571ef31dbe3a']
end

# Whens

When(/^I run the command to analyze a project$/) do
  team = 'test-team'

  @output = `./test/analyze.sh #{@project_id} #{team} #{@branch}`.chomp
end

When(/^I run the command to create a project$/) do
  team = 'test-team'
  project_name = 'java-lew'
  @branch = "master-#{Time.now.to_i.to_s + rand(1000).to_s}"
  source = 'git@github.com:ion-channel/java-lew.git'
  description = 'Java Lew'

  @output = `ion-connect project create-project --team-id #{team} --ruleset-id #{@ruleset_id} --active --branch #{@branch} #{project_name} "#{source}" "#{description}"`.chomp
end

When(/^I run the command to create a ruleset$/) do
  team = 'test-team'
  ruleset_name = 'test-ruleset'
  description = 'this is a ruleset'

  @output = `ion-connect ruleset create-ruleset --team-id #{team} #{ruleset_name} "#{description}" '#{@rules}'`.chomp
end

# Thens

Then(/^I see a response showing the project is analyzed$/) do
  expect(@output).to include('Finished about_yml scan for java-lew, valid .about.yml found.')
  expect(@output).to include('Compliance analysis completed successfully')
  expect(@output).to include('is compliant!')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end

Then(/^I see a response showing the project is created$/) do
  expect(@output).to include('"active": true')
  expect(@output).to include('"branch": "master-')
  expect(@output).to include('"source": "git@github.com:ion-channel/java-lew.git"')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end

Then(/^I see a response showing the ruleset is created$/) do
  expect(@output).to include('rules')
  expect(@output).to include('"name": "test-ruleset"')
  expect(@output).to include('"description": "The project source is required to include a valid .about.yml file.",')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end
