require 'sprockets'
require 'haml'
require 'tilt/haml'
require 'json'

Sprockets::Engines
Sprockets.register_engine '.haml', Tilt::HamlTemplate

class SprocketsEnv
  attr_reader :environment

  def initialize(opts={})
    @opts = default_opts.merge(opts)
    @environment = Sprockets::Environment.new
    @opts[:asset_paths].each do |p|
      environment.append_path p
    end
    if opts[:disable_compression]
      require File.expand_path(__FILE__, '../disable_compression')
    end
  end

  def assets_lister
    Proc.new do |env|
      assets = {}
      environment.each_logical_path {|p| assets[p] = environment.find_asset(p).digest_path}
      [
        200,
        {'Content-Type' => 'application/json', 'Cache-Control' => 'no-cache'},
        [assets.to_json]
      ]
    end
  end

  def default_opts
    {
      :asset_paths => ["assets"]
    }
  end

  def rake_task
    require 'rake/sprocketstask'
    manifest = Sprockets::Manifest.new(environment, "./public/manifest.json")
    Rake::SprocketsTask.new do |t|
      t.environment = environment
      t.output = "./public"
      t.assets = "*"
      t.keep = 2
      t.manifest = manifest
    end
  end
end
