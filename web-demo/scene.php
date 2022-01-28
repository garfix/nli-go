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
    $json = json_encode([
        "SessionId" => $sessionId,
        "ApplicationDir" => $configPath,
        "WorkDir" => $varDir,
        "Command" => "query",
        "Query" => $query
    ]);
    $fullCommand = sprintf("echo %s | netcat localhost 3333", escapeshellarg($json));

} else if ($action == "reset") {
    $json = json_encode([
        "SessionId" => $sessionId,
        "Command" => "reset"
    ]);
    $fullCommand = sprintf("echo %s | netcat localhost 3333", escapeshellarg($json));
} else {
    die("Unknown action: " . $action);
}

exec($fullCommand, $output);

header('content-type: application/json');
echo implode("\n", $output);
