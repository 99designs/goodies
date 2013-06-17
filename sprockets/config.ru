#\ -w

require 'bundler'
Bundler.setup :development
require File.expand_path('..', __FILE__) + "/environment"

env = SprocketsEnv.new()

map '/assets' do
  run env.assets_lister
end

map '/' do
  run env.environment
end
