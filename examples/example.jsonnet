local lib = import 'lib.libsonnet';
{
  who: 'world',
  say: 'hello %(who)s' % (self),
  sayAgain: 'hello %(data)s' % (lib),
}