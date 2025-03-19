var lat_lng = new google.maps.LatLng(26.218592, 127.685959);
var Options = {
  zoom: 15,      //地図の縮尺値
  center: lat_lng,    //地図の中心座標
  mapTypeId: 'roadmap'   //地図の種類
};
var map = new google.maps.Map(document.getElementById('map'), Options);
var map2 = new google.maps.Map(document.getElementById('map2'), Options);

// var lat_lng = new google.maps.LatLng(34.676375, 135.503392);

var marker = null;

// クリックしたら関数を実行
map.addListener('click', function (e) {
  getLatLng(e.latLng, map2);
  getOurLatLng(e.latLng, map);
});

function getLatLng(lat_lng, map) {

  // 座標を表示
  var lat = document.getElementById('lat').textContent = lat_lng.lat();
  var lng = document.getElementById('lng').textContent = lat_lng.lng();


  // TODO : 削除予定（緯度経度取得できてるか確認用）
  console.log({ lat });
  console.log({ lng });

  //マーカーが存在していたら削除
  if (marker) {
    marker.setMap(null);
  }

  // マーカーを再設置
  marker = new google.maps.Marker({
    position: lat_lng,
    map: map
  });
}

function getOurLatLng(lat_lng, map2) {
  // // 座標を表示
  // var lat = document.getElementById('lat').textContent = lat_lng.lat();
  // var lng = document.getElementById('lng').textContent = lat_lng.lng();
  // マーカーを設置
  marker = new google.maps.Marker({
    position: lat_lng,
    map: map2
  });
}

function sendLatLng() {
  fetch('http://localhost:8080/add-marker', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      // TODO: 削除予定（値確認用）
      // lat: 35.675012,
      // lng: 139.694388,
      lat: parseFloat(document.getElementById('lat').textContent),
      lng: parseFloat(document.getElementById('lng').textContent),
    }),
  })
  // .then((res) => {
  //   return res.json();
  // })
  // .then((data) => {
  //   console.log(data);
  // })
  // .catch((err) => {
  //   console.log(err);
  // });
}
function showAlert() {
  alert("保存されました");
  // location.reload();
}

function get() {
  fetch('http://localhost:8080/get-markers', {
    method: 'POST',
    // headers: {
    //   'Content-Type': 'application/json',
    // },
    // body: JSON.stringify({
    //   lat: parseFloat(document.getElementById('lat').textContent),
    //   lng: parseFloat(document.getElementById('lng').textContent),
    // }),
  })
}