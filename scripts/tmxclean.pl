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

sub fix_object {
  my ($el) = @_;
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Sprite') {
    my %prop = props $el;
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
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Switch') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 290);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Riser') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 291);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'SwitchableSprite') {
    my %prop = props $el;
    if (!defined $prop{image}) {
      if (($prop{initial_state} // '') ne 'false') {
        $el->removeChildNodes();
        $el->removeAttribute('type');
        $el->setAttribute('gid', 292);
      } else {
        $el->removeChildNodes();
        $el->removeAttribute('type');
        $el->setAttribute('gid', 293);
      }
      for my $node($el->getElementsByTagName('property')) {
        if ($node->getAttribute('name') eq 'initial_state') {
          $node->parent->removeChildNode($node);
        }
      }
    }
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'TnihSign') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 294);
  }
}

my $dom = XML::LibXML->load_xml(location => '../assets/maps/level.tmx');
my $doc = $dom->documentElement();
for my $el($doc->getElementsByTagName('object')) {
  fix_object($el);
}
$dom->toFile('../assets/maps/level.tmx');
