# Givens

# Whens

When(/^I run the command to analyze a project$/) do
  team = 'test-team'

  @output = `./test/analyze.sh #{@project_id} #{team} #{@branch}`.chomp
end

# Thens

Then(/^I see a response showing the project is analyzed$/) do
  expect(@output).to include("Begining compliance analysis of project (#{@project_id})")
  expect(@output).to include('Analysis requested the id is')
  expect(@output).to include('All project scans have finished')
  expect(@output).to include('Evaluating analysis for compliance')
  expect(@output).to include('Compliance analysis failed, your project is not compliant')

  # Only so the `Given previous output` step works. Once those are gone, remove this.
  $output = @output
end
