use strict;
use warnings;
use Data::Dumper;
use XML::LibXML;

sub props {
  my ($el) = @_;
  my %prop;
  for my $prop ($el->getElementsByTagName('property')) {
    my $name = $prop->getAttribute('name');
    my $value = $prop->getAttribute('value');
    $prop{$name} = $value;
  }
  return %prop;
}

sub fix_sprite {
  my ($el) = @_;
  my %prop = props $el;
  # Make clips visible in the editor.
  if ($prop{image} eq 'playerclip.png') {
    $el->removeChildNodes();
    $el->removeAttribute('type');
    $el->setAttribute('gid', 283);
  }
  if ($prop{image} eq 'objectclip.png') {
    $el->removeChildNodes();
    $el->removeAttribute('type');
    $el->setAttribute('gid', 284);
  }
}

sub fix_object {
  my ($el) = @_;
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Sprite') {
    fix_sprite($el);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'RiserFsck') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 285);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'OneWay') {
    my %prop = props $el;
    my $orientation = $prop{orientation} // 'ES';
    $el->removeChildNodes();
    $el->removeAttribute('type');
    $el->setAttribute('gid', 286) if $orientation =~ /^E/;
    $el->setAttribute('gid', 287) if $orientation =~ /^N/;
    $el->setAttribute('gid', 288) if $orientation =~ /^W/;
    $el->setAttribute('gid', 289) if $orientation =~ /^S/;
  }
}

my $dom = XML::LibXML->load_xml(location => '../assets/maps/level.tmx');
my $doc = $dom->documentElement();
for my $el($doc->getElementsByTagName('object')) {
  fix_object($el);
}
$dom->toFile('../assets/maps/level.tmx');
