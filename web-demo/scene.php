<?php

session_start();
$sessionId = session_id();

$action = $_REQUEST['action'];

$app = 'blocks';
$command = __DIR__ . '/../bin/nli';
$configPath = __DIR__ . '/../resources/' . $app;
$varDir = __DIR__ . '/../var';

if ($action == "state") {
    $query = "dom:at(E, X, Z, Y) dom:type(E, Type) dom:color(E, Color) dom:size(E, Width, Length, Height)";
    $fullCommand = sprintf('%s query -s %s -a %s -o %s "%s"', $command, $sessionId, $configPath, $varDir, $query);
} else if ($action == "reset") {
    $fullCommand = sprintf('%s reset -s %s -a %s -o %s', $command, $sessionId, $configPath, $varDir);
    $output = json_encode([]);
} else {
    die("Unknown action: " . $action);
}

exec($fullCommand, $output);

header('content-type: application/json');
echo implode("\n", $output);
