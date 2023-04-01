function(vars)
  std.lines(["%s=%s" % [k, std.escapeStringBash(vars[k])] for k in std.objectFields(vars)])