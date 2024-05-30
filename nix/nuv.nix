# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#	http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
{ lib
, stdenv
, pkgs
, fetchFromGitHub
, fetchurl
, buildGoModule
, makeWrapper
, breakpointHook
, jq
, curl
, kubectl 
, eksctl 
, kind 
, k3sup 
, coreutils
}:

let 
   branch = "3.0.0";
   version = "3.0.1-beta.2405292059";
in buildGoModule rec {

  pname = "nuv";

  inherit branch version;

  nativeBuildInputs = [ makeWrapper jq curl breakpointHook] ;

  buildInputs = [ kubectl eksctl kind k3sup coreutils ];

  subPackages = ["."];

  src = fetchFromGitHub {
    owner = "nuvolaris";
    repo = "nuv";
    rev = version;
    sha256 = "sha256-MdnBvlA4S2Mi/bcbE+O02x+wvlIrsK1Zc0dySz4FB/w=";
  };
  
  vendorHash =  "sha256-JkQbQ2NEaumXbAfsv0fNiQf/EwMs3SDLHvu7c/bU7fU=";

  doCheck = false; 

  ldFlags =  [
    "-X main.NuvVersion=${version}"
     "-X main.NuvBranch=${branch}"
  ];

  meta = with lib; {
    description = "Nuvolaris Almighty CLI tool";
    license = licenses.asl20;
    homepage = "https://nuvolaris.io/";
    maintainers = with maintainers; [ aacebedo ];
    mainProgram = "nuv";
  };

  postInstall  = ''
     makeWrapper ${coreutils}/bin/coreutils $out/bin/coreutils
     makeWrapper ${kubectl}/bin/kubectl $out/bin/kubectl
     makeWrapper ${eksctl}/bin/eksctl $out/bin/eksctl
     makeWrapper ${kind}/bin/kind $out/bin/kind
     makeWrapper ${k3sup}/bin/k3sup $out/bin/k3sup
  '';


}