require 'godutch'

module TestGoDutch
  include GoDutch::Reactor

  def check_test
    success("Everything is o'right.")
    metric({ 'okay' => 1 })
    return 'check_test output'
  end

  def check_second_test
    critical("Here stuff is getting hard!")
    return 'something else'
  end
end

GoDutch.run(TestGoDutch)
