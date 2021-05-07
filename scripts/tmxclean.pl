use strict;
use warnings;
use Data::Dumper;
use XML::LibXML;

$Data::Dumper::Sortkeys = 1;
$Data::Dumper::Useqq = 1;

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

sub remove_props {
  my ($el, @props) = @_;
  my %props = map { $_ => 1 } @props;
  for my $node($el->getElementsByTagName('property')) {
    if (exists $props{$node->getAttribute('name')}) {
      $node->parentNode->removeChild($node);
    }
  }
}

my %objects = ();

sub fix_object {
  my ($el) = @_;
  $el->hasAttribute('type') or return;
  my $type = $el->getAttribute('type');
  ++$objects{$type};
  my %prop = props $el;
  if ($type eq 'Sprite') {
    my $img = $prop{image};
    ++$objects{'Sprite=' . $img};
    if ($img eq 'playerclip.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 283);
    }
    if ($img eq 'objectclip.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 284);
    }
    if ($img eq 'gradient_left_right.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 298);
    }
    if ($img eq 'gradient_top_bottom.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 299);
    }
    if ($img eq 'gradient_outside_inside.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 300);
    }
    if ($img eq 'arrow32.png') {
      remove_props $el, 'image', 'orientation';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 302) if $prop{orientation} =~ /EN|NE/;
      $el->setAttribute('gid', 304) if $prop{orientation} =~ /NW|WN/;
      $el->setAttribute('gid', 306) if $prop{orientation} =~ /SW|WS/;
      $el->setAttribute('gid', 308) if $prop{orientation} =~ /SE|ES/;
    }
  }
  if ($type eq 'RiserFsck') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 285);
  }
  if ($type eq 'OneWay') {
    my $orientation = $prop{orientation} // 'ES';
    remove_props $el, 'orientation';
    $el->removeAttribute('type');
    $el->setAttribute('gid', 286) if $orientation =~ /^E/;
    $el->setAttribute('gid', 287) if $orientation =~ /^N/;
    $el->setAttribute('gid', 288) if $orientation =~ /^W/;
    $el->setAttribute('gid', 289) if $orientation =~ /^S/;
  }
  if ($type eq 'Switch') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 290);
  }
  if ($type eq 'Riser') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 291);
  }
  if ($type eq 'SwitchableSprite') {
    if (!defined $prop{image} && !$el->hasAttribute('gid')) {
      if (($prop{invert} // '') eq 'true') {
        $el->removeAttribute('type');
        $el->setAttribute('gid', 292);
      } else {
        $el->removeAttribute('type');
        $el->setAttribute('gid', 293);
      }
      remove_props $el, 'invert';
    }
  }
  if ($type eq 'TnihSign') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 294);
  }
  if ($type eq 'QuestionBlock') {
    remove_props $el, 'kaizo';
    $el->removeAttribute('type');
    $el->setAttribute('gid', (($prop{kaizo} // '') eq 'true') ? 296 : 295);
  }
  if ($type eq 'AppearBlock') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 297);
  }
}

my $dom = XML::LibXML->load_xml(location => '../assets/maps/level.tmx');
my $doc = $dom->documentElement();
for my $el($doc->getElementsByTagName('object')) {
  fix_object($el);
}
$dom->toFile('../assets/maps/level.tmx');
print Dumper \%objects;
