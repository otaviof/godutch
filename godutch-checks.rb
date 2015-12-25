require 'godutch'

module TestGoDutch
  include GoDutch::Reactor
  extend self

  def check_test
    success("Everything is o'right.")
    metric({ 'okay' => 1 })
    return 'check_test output'
  end

  def dummy_method
    puts 'I should never be called.'
  end

  def check_second_test
    puts 'Foo'
  end
end

GoDutch.run(TestGoDutch)
