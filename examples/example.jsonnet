local lib = import 'lib.libsonnet';
{
  who: 'world',
  say: 'hello %(who)s' % (self),
  sayAgain: 'hello %(data)s' % (lib),
  sayExtStr: 'hello %s' % [std.extVar('hello')],
  calcExtCode: 'two + two = %d' % [std.extVar('calc')],
}