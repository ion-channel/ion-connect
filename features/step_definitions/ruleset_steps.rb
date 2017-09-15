# Givens

Given(/^I have a set of rules$/) do
  @rules = ['c30b917956c3040daa2c571ef31dbe3a']
end

Given(/^I have a ruleset id$/) do
  steps %Q{
    Given I have a set of rules
    When I run the command to create a ruleset
  }

  json = JSON.parse(@output)
  @ruleset_id = json['id']
end

# Whens

When(/^I run the command to create a ruleset$/) do
  team = 'test-team'
  ruleset_name = 'test-ruleset'
  description = 'this is a ruleset'

  @output = `ion-connect ruleset create-ruleset --team-id #{team} #{ruleset_name} "#{description}" '#{@rules}'`.chomp
end

# Thens

Then(/^I see a response showing the ruleset is created$/) do
  expect(@output).to include('rules')
  expect(@output).to include('"name": "test-ruleset"')
  expect(@output).to include('"description": "The project source is required to include a valid .about.yml file.",')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end
