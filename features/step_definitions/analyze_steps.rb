# Givens

# Whens

When(/^I run the command to analyze a project$/) do
  team = 'test-team'

  @output = `./test/analyze.sh #{@project_id} #{team} #{@branch}`.chomp
end

# Thens

Then(/^I see a response showing the project is analyzed$/) do
  expect(@output).to include('Finished about_yml scan for java-lew, valid .about.yml found.')
  expect(@output).to include('Compliance analysis completed successfully')
  expect(@output).to include('is compliant!')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end
