import sys
from unittest.mock import MagicMock

# Mock mcp module
sys.modules["mcp"] = MagicMock()
sys.modules["mcp.server"] = MagicMock()

# Now import server
from server import CrossplaneGenerator

def test_generate():
    gen = CrossplaneGenerator()
    
    # Test Instance
    config_instance = {"name": "test-vm", "region": "us-west-2", "instanceType": "m5.large"}
    result_instance = gen.generate("aws", "instance", config_instance)
    print("--- Instance AWS ---")
    print(result_instance)
    assert "kind: ComputeVM" in result_instance
    assert "provider: aws" in result_instance
    assert "region: us-west-2" in result_instance
    
    # Test Database
    config_db = {"name": "test-db", "region": "us-central1", "engine": "mysql"}
    result_db = gen.generate("gcp", "database", config_db)
    print("\n--- Database GCP ---")
    print(result_db)
    assert "kind: DatabaseInstance" in result_db
    assert "provider: gcp" in result_db
    assert "engine: mysql" in result_db

    print("\n✅ Verification Successful!")

if __name__ == "__main__":
    test_generate()
