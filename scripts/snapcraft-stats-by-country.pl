#!/usr/bin/perl

use strict;
use warnings;
use POSIX qw(strftime);
use JSON qw(decode_json);

my $top = 15;

my $now = time;
my $start_date = strftime '%Y-%m-%d', gmtime $now - 86400 * 90;
my $end_date   = strftime '%Y-%m-%d', gmtime $now - 86400 * 2;

open my $json, '-|', qw(snap run snapcraft metrics aaaaxy),
	'--name' => 'weekly_installed_base_by_country',
	'--format' => 'json',
	'--start' => $start_date,
	'--end' => $end_date
	or die "open: $!";
my $raw = do { undef local $/; <$json>; }
	or die "read: $!";
close $json
	or die "close: $!";
my $data = decode_json $raw;

# Turn into a hashmap.
my %sheet;
for my $series (@{$data->{series}}) {
	my $key = $series->{name};
	next if $key eq 'none';
	my $i = 0;
	for my $value(@{$series->{values}}) {
		my $bucket = $data->{buckets}[$i];
		$sheet{$key}{$bucket} = $value;
		++$i;
	}
}

# Do top N.
my @buckets = sort @{$data->{buckets}};
my $newest_bucket = [@buckets]->[-1];
my @keys = sort {
		($sheet{$b}{$newest_bucket} // 0) <=> ($sheet{$a}{$newest_bucket} // 0)
	} keys %sheet;
my @topkeys = @keys[0..$top - 1];

# Compute the "other" bucket.
my %other;
my %is_topkey = map { $_ => 1 } @topkeys;
while (my ($key, $row) = each %sheet) {
	next if $is_topkey{$key};
	while (my ($bucket, $value) = each %$row) {
		next unless defined $value;
		$other{$bucket} += $value;
	}
}

my @rows = (
	(map {
		[$_, $sheet{$_}]
	} @topkeys),
	['other', \%other]
);

my $header = "set timefmt '%Y-%m-%d'\nset xdata time\nplot";
my $header_comma = ' ';
my $body = '';
for my $item(@rows) {
	my ($key, $row) = @$item;
	$header .= "$header_comma'-' using 1:2 with linespoints title '$key'";
	$header_comma = ', ';
	for my $bucket(@buckets) {
		$body .= "$bucket @{[$row->{$bucket} // 0]}\n";
	}
	$body .= "e\n";
}

print "$header\n$body";
