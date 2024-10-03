  #!/bin/bash
  #
  # Copyright IBM Corp All Rights Reserved
  #
  # SPDX-License-Identifier: Apache-2.0
  #

  # This is a collection of bash functions used by different scripts

  # imports
  . scripts/utils.sh

  # Enable TLS
  export CORE_PEER_TLS_ENABLED=true

  # Path to orderer CA files for TLS (for 3 orderers)
  export ORDERER1_CA=${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
  export ORDERER2_CA=${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
  export ORDERER3_CA=${PWD}/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem

  # Path to peer CAs for each organization (with 2 peers per org)
  export PEER0_PRODUCER_CA=${PWD}/organizations/peerOrganizations/producer.example.com/tlsca/tlsca.producer.example.com-cert.pem
  export PEER1_PRODUCER_CA=${PWD}/organizations/peerOrganizations/producer.example.com/tlsca/tlsca.producer.example.com-cert.pem
  export PEER0_SUPPLIER_CA=${PWD}/organizations/peerOrganizations/supplier.example.com/tlsca/tlsca.supplier.example.com-cert.pem
  export PEER1_SUPPLIER_CA=${PWD}/organizations/peerOrganizations/supplier.example.com/tlsca/tlsca.supplier.example.com-cert.pem
  export PEER0_wholesaler_CA=${PWD}/organizations/peerOrganizations/wholesaler.example.com/tlsca/tlsca.wholesaler.example.com-cert.pem
  export PEER1_wholesaler_CA=${PWD}/organizations/peerOrganizations/wholesaler.example.com/tlsca/tlsca.wholesaler.example.com-cert.pem

  # Path to orderer TLS certificates (adjust paths if necessary)
  export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.crt
  export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.key

  # Set environment variables for the peer org
  setGlobals() {
    local USING_ORG=""
    local PEER_ID=""
    if [ -z "$OVERRIDE_ORG" ]; then
      USING_ORG=$1
      PEER_ID=$2
    else
      USING_ORG="${OVERRIDE_ORG}"
      PEER_ID="${OVERRIDE_PEER}"
    fi
    infoln "Using organization ${USING_ORG} and peer ${PEER_ID}"
    
    if [ $USING_ORG -eq 1 ]; then
      export CORE_PEER_LOCALMSPID="ProducerMSP"
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_PRODUCER_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/producer.example.com/users/Admin@producer.example.com/msp
        export CORE_PEER_ADDRESS=localhost:7051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_PRODUCER_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/producer.example.com/users/Admin@producer.example.com/msp
        export CORE_PEER_ADDRESS=localhost:8051
      fi

    elif [ $USING_ORG -eq 2 ]; then
      export CORE_PEER_LOCALMSPID="SupplierMSP"
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_SUPPLIER_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier.example.com/users/Admin@supplier.example.com/msp
        export CORE_PEER_ADDRESS=localhost:9051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_SUPPLIER_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/supplier.example.com/users/Admin@supplier.example.com/msp
        export CORE_PEER_ADDRESS=localhost:10051
      fi

    elif [ $USING_ORG -eq 3 ]; then
      export CORE_PEER_LOCALMSPID="WholesalerMSP"
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_wholesaler_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/wholesaler.example.com/users/Admin@wholesaler.example.com/msp
        export CORE_PEER_ADDRESS=localhost:11051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_TLS_ROOTCERT_FILE=$PEER1_wholesaler_CA
        export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/wholesaler.example.com/users/Admin@wholesaler.example.com/msp
        export CORE_PEER_ADDRESS=localhost:12051
      fi

    else
      errorln "ORG Unknown"
    fi

    if [ "$VERBOSE" == "true" ]; then
      env | grep CORE
    fi
  }

  # Set environment variables for use in the CLI container
  setGlobalsCLI() {
    setGlobals $1 $2

    local USING_ORG=""
    local PEER_ID=""
    if [ -z "$OVERRIDE_ORG" ]; then
      USING_ORG=$1
      PEER_ID=$2
    else
      USING_ORG="${OVERRIDE_ORG}"
      PEER_ID="${OVERRIDE_PEER}"
    fi

    if [ $USING_ORG -eq 1 ]; then
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_ADDRESS=peer0.producer.example.com:7051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_ADDRESS=peer1.producer.example.com:8051
      fi

    elif [ $USING_ORG -eq 2 ]; then
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_ADDRESS=peer0.supplier.example.com:9051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_ADDRESS=peer1.supplier.example.com:10051
      fi

    elif [ $USING_ORG -eq 3 ]; then
      if [ $PEER_ID -eq 0 ]; then
        export CORE_PEER_ADDRESS=peer0.wholesaler.example.com:11051
      elif [ $PEER_ID -eq 1 ]; then
        export CORE_PEER_ADDRESS=peer1.wholesaler.example.com:12051
      fi

    else
      errorln "ORG Unknown"
    fi
  }

  # parsePeerConnectionParameters $@
  # Helper function that sets the peer connection parameters for a chaincode operation
  parsePeerConnectionParameters() {
    PEER_CONN_PARMS=()
    PEERS=""
    while [ "$#" -gt 0 ]; do
      setGlobals $1 $2
      PEER="peer$2.org$1"
      ## Set peer addresses
      if [ -z "$PEERS" ]; then
        PEERS="$PEER"
      else
        PEERS="$PEERS $PEER"
      fi
      PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" --peerAddresses $CORE_PEER_ADDRESS)
      ## Set path to TLS certificate
      CA_VAR="PEER${2}_ORG${1}_CA"
      TLSINFO=(--tlsRootCertFiles "${!CA_VAR}")
      PEER_CONN_PARMS=("${PEER_CONN_PARMS[@]}" "${TLSINFO[@]}")
      # Shift by two to get to the next organization-peer pair
      shift 2
    done
  }

  verifyResult() {
    if [ $1 -ne 0 ]; then
      fatalln "$2"
    fi
  }
