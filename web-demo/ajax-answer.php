<?php

$request = isset($_REQUEST['request']) ? $_REQUEST['request'] : null;
$app = isset($_REQUEST['app']) ? $_REQUEST['app'] : "dbpedia";

if (!in_array($app, ['dbpedia', 'blocks'])) {
    die('Unsupported app');
}
session_start();
$sessionId = session_id();

$command = __DIR__ . '/../bin/nli';
$configPath = __DIR__ . '/../resources/' . $app;
$varDir = __DIR__ . '/../var';

$start = microtime(true);

$message = json_decode($request);
$json = json_encode([
    "SessionId" => $sessionId,
    "ApplicationDir" => $configPath,
    "WorkDir" => $varDir,
    "Command" => "send",
    "Message" => $message,
]);
$fullCommand = sprintf("echo %s | netcat localhost 3333", escapeshellarg($json));
exec($fullCommand, $output);

$end = microtime(true);
$duration = sprintf("%.2f", $end - $start);

$result = json_decode(implode($output), true);
$response = json_encode($result, JSON_PRETTY_PRINT);
if ($result === null) {
    $response = "ERROR executing\n" . $fullCommand . "\n" . implode($output);
}

logResult($request, $result, $duration, $sessionId);

header('content-type: application/json');
echo $response;

function logResult($request, $result, $duration, $sessionId)
{
    $handle = fopen(__DIR__ . '/../var/log/' . date('Y-m') . '-queries.log', 'a');

    if ($result) {
        unset($result['Productions']);
    }

    $answer = json_encode($result);
    $answer = strlen($answer) < 300 ? $answer : substr($answer, 0, 300) . "...";

    fwrite($handle, "#  " . date('Y-m-d H:i:s') . " " . $sessionId .  " (" . $duration . "s)\n");
    fwrite($handle, "Q: " . $request . "\n");
    fwrite($handle, "A: " . $answer . "\n");

    fwrite($handle, "\n");
    fclose($handle);
}