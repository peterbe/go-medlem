#!/bin/bash


echo "GET"
curl "http://localhost:5000/staff?email=pbengtsson@mozilla.com&email=two@example.com"
echo
echo "----"

echo "GET (invalid email)"
curl "http://localhost:5000/staff?email=oneexamplecom"
echo
echo "----"

echo "GET (empty)"
curl "http://localhost:5000/staff"
echo
echo "----"

echo "POST FORM DATA"
curl -X POST -F "email=one@example.com" -F "email=pbengtsson@mozilla.com" http://localhost:5000/staff
echo
echo "----"

echo "POST JSON DATA"
curl -X POST -H "Content-Type: application/json" \
 -d '{"email": ["one@example.com", "pbengtsson@mozilla.com"]}' http://localhost:5000/staff
echo
echo "----"
