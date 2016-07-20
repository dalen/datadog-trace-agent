
def go_build(program, opts={})
  default_cmd = "go build -a"
  if ENV["INCREMENTAL_BUILD"] then
    default_cmd = "go build -i"
  end
  opts = {
    :cmd => default_cmd
  }.merge(opts)

  dd = 'main'
  commit = `git rev-parse --short HEAD`.strip
  branch = `git rev-parse --abbrev-ref HEAD`.strip
  date = `date`.strip
  goversion = `go version`.strip
  ldflags = {
    "#{dd}.BuildDate" => "#{date}",
    "#{dd}.GitCommit" => "#{commit}",
    "#{dd}.GitBranch" => "#{branch}",
    "#{dd}.GoVersion" => "#{goversion}",
  }.map do |k,v|
    if goversion.include?("1.4")
      "-X #{k} '#{v}'"
    else
      "-X '#{k}=#{v}'"
    end
  end.join ' '

  cmd = opts[:cmd]
  sh "#{cmd} -ldflags \"#{ldflags}\" #{program}"
end


def go_lint(path)
  out = `golint #{path}/*.go`
  errors = out.split("\n")
  puts "#{errors.length} linting issues found"
  if errors.length > 0
    puts out
    fail
  end
end

def go_vet(path)
  sh "go vet #{path}"
end

def go_test(path, opts={})
  opts = {
    :v => false,
    :include => "raclette"
  }.merge(opts)

  paths = [path]
  if opts[:include]
    deps = `go list -f '{{ join .Deps "\\n"}}' #{path} | sort | uniq`.split("\n").select do |p|
      p.include? opts[:include]
    end
    paths = paths.concat(deps)
  end

  v = opts[:v] ? "-v" : ""
  sh "go test #{v} #{paths.join(' ')}"
end

# return the dependencies of all the packages who start with the root path
def go_pkg_deps(pkgs, root_path)
  deps = []
  pkgs.each do |pkg|
    deps << pkg
    `go list -f '{{ join .Deps "\\n"}}' #{pkg}`.split("\n").select do |path|
      if path.start_with? root_path
        deps << path
      end
    end
  end
  return deps.sort.uniq
end

def go_fmt(path)
  out = `go fmt ./...`
  errors = out.split("\n")
  if errors.length > 0
    puts out
    fail
  end
end

