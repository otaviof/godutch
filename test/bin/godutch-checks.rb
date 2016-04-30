#!/usr/bin/env ruby

require 'godutch'

# Simple GoDutch example
module TestGoDutch
  include GoDutch::Reactor

  def check_test
    success('Everything is o\'right.')
    metric('okay' => 1)
    'check_test output'
  end

  def check_second_test
    critical('Here stuff is getting hard!')
    'something else'
  end

  def check_third_test
  end

  def check_forth_test
  end
end

GoDutch.run(TestGoDutch)
