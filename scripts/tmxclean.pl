use strict;
use warnings;
use Data::Dumper;
use XML::LibXML;

sub fix_sprite {
  my ($el) = @_;
  my %prop;
  for my $prop ($el->getElementsByTagName('property')) {
    my $name = $prop->getAttribute('name');
    my $value = $prop->getAttribute('value');
    $prop{$name} = $value;
  }
  # Make clips visible in the editor.
  if ($prop{image} eq 'playerclip.png') {
    $el->removeChildNodes();
    $el->removeAttribute('type');
    $el->setAttribute('gid', 1000);
  }
  if ($prop{image} eq 'objectclip.png') {
    $el->removeChildNodes();
    $el->removeAttribute('type');
    $el->setAttribute('gid', 1001);
  }
}

sub fix_object {
  my ($el) = @_;
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Sprite') {
    fix_sprite($el);
  }
  # TODO: Also show one-ways.
}

my $dom = XML::LibXML->load_xml(location => '../assets/maps/level.tmx');
my $doc = $dom->documentElement();
for my $el($doc->getElementsByTagName('object')) {
  fix_object($el);
}
$dom->toFile('../assets/maps/level.tmx');
