<?php

$action = $_REQUEST['action'];

$app = 'blocks';
$command = realpath(__DIR__ . '/../bin/nli');
$configPath = realpath(__DIR__ . '/../resources/' . $app);
$varDir = realpath(__DIR__ . '/../var');

session_start();
$sessionId = $app . "_" . session_id();

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
