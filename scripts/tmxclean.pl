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

sub remove_props {
  my ($el, @props) = @_;
  my %props = map { $_ => 1 } @props;
  for my $node($el->getElementsByTagName('property')) {
    if (exists $props{$node->getAttribute('name')}) {
      $node->parentNode->removeChild($node);
    }
  }
}

sub fix_object {
  my ($el) = @_;
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'Sprite') {
    my %prop = props $el;
    if ($prop{image} eq 'playerclip.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 283);
    }
    if ($prop{image} eq 'objectclip.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 284);
    }
    if ($prop{image} eq 'gradient_left_right.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 298);
    }
    if ($prop{image} eq 'gradient_top_bottom.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 299);
    }
    if ($prop{image} eq 'gradient_outside_inside.png') {
      remove_props $el, 'image';
      $el->removeAttribute('type');
      $el->setAttribute('gid', 300);
    }
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'RiserFsck') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 285);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'OneWay') {
    my %prop = props $el;
    my $orientation = $prop{orientation} // 'ES';
    remove_props $el, 'orientation';
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
        $el->removeAttribute('type');
        $el->setAttribute('gid', 292);
      } else {
        $el->removeAttribute('type');
        $el->setAttribute('gid', 293);
      }
      remove_props $el, 'initial_state';
    }
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'TnihSign') {
    $el->removeAttribute('type');
    $el->setAttribute('gid', 294);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'QuestionBlock') {
    my %prop = props $el;
    remove_props $el, 'kaizo';
    $el->removeAttribute('type');
    $el->setAttribute('gid', (($prop{kaizo} // '') eq 'true') ? 296 : 295);
  }
  if ($el->hasAttribute('type') && $el->getAttribute('type') eq 'AppearBlock') {
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
