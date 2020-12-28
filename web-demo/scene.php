<?php

session_start();
$sessionId = session_id();

$app = 'blocks';
$command = __DIR__ . '/../bin/nli';
$configPath = __DIR__ . '/../resources/' . $app;
$varDir = __DIR__ . '/../var';

$query = "dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)";
$fullCommand = sprintf('%s -s %s -c %s -d %s -r json -q "%s"', $command, $sessionId, $configPath, $varDir, $query);
exec($fullCommand, $output);

header('content-type: application/json');
echo implode("\n", $output);
